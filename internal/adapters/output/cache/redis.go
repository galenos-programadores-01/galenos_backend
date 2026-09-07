package cache

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache define la interfaz para operaciones de caché de alto rendimiento.
type Cache interface {
	Get(ctx context.Context, key string, target any) bool
	Set(ctx context.Context, key string, value any, ttl time.Duration)
	Delete(ctx context.Context, keys ...string)
	IsAvailable() bool
}

// RedisCache implementa Cache utilizando Redis v9 con degradación grácil.
type RedisCache struct {
	client     *redis.Client
	defaultTTL time.Duration
	available  bool
	mu         sync.RWMutex
}

// NewRedisCache inicializa el cliente Redis y verifica la conectividad sin bloquear el arranque.
func NewRedisCache(addr, password string, db int, defaultTTL time.Duration) *RedisCache {
	rc := &RedisCache{
		defaultTTL: defaultTTL,
	}

	if addr == "" {
		return rc
	}

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  1 * time.Second,
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("[REDIS] No se pudo conectar a Redis en %s (%v). Modo fallback a BD activo.", addr, err)
		rc.available = false
	} else {
		log.Printf("[REDIS] Conexión establecida con éxito en %s (DB %d).", addr, db)
		rc.available = true
	}

	rc.client = client
	return rc
}

// IsAvailable retorna si Redis está respondiendo activamente.
func (r *RedisCache) IsAvailable() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.available && r.client != nil
}

// Get obtiene un valor serializado en JSON y lo deserializa en target. Retorna true si hubo HIT.
func (r *RedisCache) Get(ctx context.Context, key string, target any) bool {
	if !r.IsAvailable() {
		return false
	}

	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err != redis.Nil {
			log.Printf("[REDIS GET ERROR] key=%s: %v", key, err)
		}
		return false
	}

	if err := json.Unmarshal(val, target); err != nil {
		log.Printf("[REDIS JSON UNMARSHAL ERROR] key=%s: %v", key, err)
		return false
	}

	return true
}

// Set serializa un valor en JSON y lo almacena con su TTL correspondiente.
func (r *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) {
	if !r.IsAvailable() {
		return
	}

	if ttl <= 0 {
		ttl = r.defaultTTL
	}

	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("[REDIS JSON MARSHAL ERROR] key=%s: %v", key, err)
		return
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		log.Printf("[REDIS SET ERROR] key=%s: %v", key, err)
	}
}

// Delete elimina una o más claves del caché.
func (r *RedisCache) Delete(ctx context.Context, keys ...string) {
	if !r.IsAvailable() || len(keys) == 0 {
		return
	}
	_ = r.client.Del(ctx, keys...).Err()
}
