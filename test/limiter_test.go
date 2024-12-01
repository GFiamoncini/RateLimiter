package limiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/GFiamoncini/RateLimiter/limiter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRedis struct {
	mock.Mock
}

func (m *MockRedis) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockRedis) Set(key string, value string, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

type MockRateLimiter struct {
	mock.Mock
}

func (m *MockRateLimiter) Allow(ctx context.Context, key string, limit int) (bool, time.Duration) {
	args := m.Called(ctx, key, limit)
	return args.Bool(0), args.Get(1).(time.Duration)
}

// Teste: Limitação por IP
func TestRateLimiter_LimitByIP(t *testing.T) {
	mockRateLimiter := new(MockRateLimiter)

	// Simulando que o IP "192.168.1.1" já fez 6 requisições, então ele excedeu o limite de 5.
	mockRateLimiter.On("Allow", mock.Anything, "ip:192.168.1.1", 5).Return(false, time.Minute)

	// Criando o RateLimiter com a estratégia mockada
	limiter := limiter.NewRateLimiter(mockRateLimiter)

	// Testando se o IP excede o limite
	allowed, _ := limiter.Allow(context.Background(), "ip:192.168.1.1", 5)
	assert.False(t, allowed, "O IP deve ter excedido o limite de requisições")

	// Verificando se o mock foi chamado corretamente
	mockRateLimiter.AssertExpectations(t)
}

// Teste: Limitação por Token
func TestRateLimiter_LimitByToken(t *testing.T) {
	mockRateLimiter := new(MockRateLimiter)

	// Simulando que o token "abc123" já fez 11 requisições, excedendo o limite de 10.
	mockRateLimiter.On("Allow", mock.Anything, "token:ABC", 10).Return(false, time.Minute)

	// Criando o RateLimiter com a estratégia mockada
	limiter := limiter.NewRateLimiter(mockRateLimiter)

	// Testando se o token excede o limite
	allowed, _ := limiter.Allow(context.Background(), "token:ABC", 10)
	assert.False(t, allowed, "O token deve ter excedido o limite de requisições")

	// Verificando se o mock foi chamado corretamente
	mockRateLimiter.AssertExpectations(t)
}

// Teste: Limitação por IP - Aceitar se dentro do limite
func TestRateLimiter_LimitByIP_Accept(t *testing.T) {
	mockRateLimiter := new(MockRateLimiter)

	// Simulando que o IP "192.168.1.1" fez 4 requisições, então está dentro do limite de 5.
	mockRateLimiter.On("Allow", mock.Anything, "ip:192.168.1.1", 5).Return(true, time.Duration(0)) // Corrigido para time.Duration(0)

	// Criando o RateLimiter com a estratégia mockada
	limiter := limiter.NewRateLimiter(mockRateLimiter)

	// Testando se o IP está dentro do limite
	allowed, _ := limiter.Allow(context.Background(), "ip:192.168.1.1", 5)
	assert.True(t, allowed, "O IP deve estar dentro do limite de requisições")

	// Verificando se o mock foi chamado corretamente
	mockRateLimiter.AssertExpectations(t)
}

// Teste: Limitação por Token - Aceitar se dentro do limite
func TestRateLimiter_LimitByToken_Accept(t *testing.T) {
	mockRateLimiter := new(MockRateLimiter)

	// Simulando que o token "abc123" fez 9 requisições, então está dentro do limite de 10.
	mockRateLimiter.On("Allow", mock.Anything, "token:ABC", 10).Return(true, time.Duration(0)) // Corrigido para time.Duration(0)

	// Criando o RateLimiter com a estratégia mockada
	limiter := limiter.NewRateLimiter(mockRateLimiter)

	// Testando se o token está dentro do limite
	allowed, _ := limiter.Allow(context.Background(), "token:ABC", 10)
	assert.True(t, allowed, "O token deve estar dentro do limite de requisições")

	// Verificando se o mock foi chamado corretamente
	mockRateLimiter.AssertExpectations(t)
}

// Teste: Sobreposição de limite de IP e Token
func TestRateLimiter_OverriddenByToken(t *testing.T) {
	mockRateLimiter := new(MockRateLimiter)

	// Simulando que o IP "192.168.1.1" fez 5 requisições, e o token "abc123" fez 11 requisições.
	mockRateLimiter.On("Allow", mock.Anything, "ip:192.168.1.1", 5).Return(false, time.Minute)
	mockRateLimiter.On("Allow", mock.Anything, "token:ABC", 10).Return(false, time.Minute)

	// Criando o RateLimiter com a estratégia mockada
	limiter := limiter.NewRateLimiter(mockRateLimiter)

	// Testando se o IP "192.168.1.1" é rejeitado (excedeu o limite)
	allowed, _ := limiter.Allow(context.Background(), "ip:192.168.1.1", 5)
	assert.False(t, allowed, "O IP deve ter excedido o limite de requisições")

	// Testando se o token "abc123" é rejeitado (excedeu o limite)
	allowed, _ = limiter.Allow(context.Background(), "token:ABC", 10)
	assert.False(t, allowed, "O token deve ter excedido o limite de requisições")

	// Verificando se o mock foi chamado corretamente para o IP e para o Token
	mockRateLimiter.AssertExpectations(t)
}
