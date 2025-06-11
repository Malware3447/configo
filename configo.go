package configo

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config interface {
	Env() string
}

type Database struct {
	Type          string        `yaml:"type" env-required:"true"`
	Host          string        `yaml:"host" env-required:"true"`
	Port          int           `yaml:"port" env-required:"true"`
	Name          string        `yaml:"name" env-required:"true"`
	User          string        `yaml:"user" env-required:"true"`
	Password      string        `yaml:"password" env-required:"true"`
	Schema        string        `yaml:"schema" env-default:"public"`
	MigrationPath string        `yaml:"migrationPath" env-required:"true"`
	MaxAttempts   int           `yaml:"maxAttempts" env-required:"true"`
	AttemptDelay  time.Duration `yaml:"attemptDelay" env-required:"true"`
}

type Redis struct {
	Host string `yaml:"host" env-required:"true"`
	Port int    `yaml:"port" env-required:"true"`
	Db   int    `yaml:"db" `
}

type KafkaProducer struct {
	Brokers      []string      `yaml:"brokers" env-separator:"," env-required:"true"`
	RequiredAcks int           `yaml:"requiredAcks" env-default:"1"` // Уровень подтверждения: 0=None, 1=Leader, -1=All
	Async        bool          `yaml:"async" env-default:"false"`
	BatchSize    int           `yaml:"batchSize" env-default:"100"`
	BatchTimeout time.Duration `yaml:"batchTimeout" env-default:"1s"`
	WriteTimeout time.Duration `yaml:"writeTimeout" env-default:"10s"`
	MaxAttempts  int           `yaml:"maxAttempts" env-default:"3"`
}

type KafkaConsumer struct {
	Brokers     []string `yaml:"brokers" env-separator:"," env-required:"true"`
	GroupID     string   `yaml:"groupId" env-required:"true"`
	Topics      []string `yaml:"topics" env-separator:"," env-required:"true"`
	StartOffset string   `yaml:"startOffset" env-default:"latest"` // 'latest' или 'earliest'

	MinBytes int           `yaml:"minBytes" env-default:"10000"`    // 10KB - Минимальный размер пакета для Fetch
	MaxBytes int           `yaml:"maxBytes" env-default:"10000000"` // 10MB - Максимальный размер пакета для Fetch
	MaxWait  time.Duration `yaml:"maxWait" env-default:"1s"`        // Макс. время ожидания MinBytes

	CommitInterval    time.Duration `yaml:"commitInterval" env-default:"1s"`    // Интервал авто-коммита (0 - отключает авто-коммит)
	HeartbeatInterval time.Duration `yaml:"heartbeatInterval" env-default:"3s"` // Частота отправки heartbeat брокеру
	SessionTimeout    time.Duration `yaml:"sessionTimeout" env-default:"30s"`   // Таймаут сессии консьюмера
	RebalanceTimeout  time.Duration `yaml:"rebalanceTimeout" env-default:"60s"` // Таймаут для ребалансировки

	DialTimeout  time.Duration `yaml:"dialTimeout" env-default:"3s"`   // Таймаут подключения к брокеру
	ReadTimeout  time.Duration `yaml:"readTimeout" env-default:"30s"`  // Таймаут чтения сообщений
	WriteTimeout time.Duration `yaml:"writeTimeout" env-default:"10s"` // Таймаут записи (для коммитов и т.д.)
	MaxAttempts  int           `yaml:"maxAttempts" env-default:"3"`    // Макс. кол-во попыток для некоторых операций
}

type KafkaTopics struct {
	List              []string `yaml:"list" env-separator:"," env-required:"true"`
	NumPartitions     int      `yaml:"numPartitions" env-required:"true"`
	ReplicationFactor int      `yaml:"replicationFactor" env-required:"true"`
}

type GrpcServer struct {
	// Адрес для прослушивания
	Host string `yaml:"host" env:"GRPC_HOST" env-default:"0.0.0.0"`        // Хост, на котором сервер будет слушать
	Port int    `yaml:"port" env:"GRPC_PORT,required" env-required:"true"` // Порт для прослушивания

	// Конфигурация TLS
	EnableTLS    bool   `yaml:"enableTLS" env:"GRPC_ENABLE_TLS" env-default:"false"` // Включить ли TLS
	CertFile     string `yaml:"certFile" env:"GRPC_CERT_FILE"`                       // Путь к файлу сертификата сервера
	KeyFile      string `yaml:"keyFile" env:"GRPC_KEY_FILE"`                         // Путь к файлу приватного ключа сервера
	ClientCAFile string `yaml:"clientCAFile" env:"GRPC_CLIENT_CA_FILE"`              // Путь к файлу CA сертификата клиента для mTLS (опционально)

	// KeepAlive параметры сервера (grpc.KeepaliveServerParameters - KASP)
	// Управляют поведением сервера относительно keep-alive от клиентов и инициируемых сервером пингов.
	// Если значение Duration равно 0, часто используется значение по умолчанию gRPC (бесконечность или очень большое).
	KeepAliveMaxConnectionIdle     time.Duration `yaml:"keepAliveMaxConnectionIdle" env:"GRPC_KEEP_ALIVE_MAX_CONNECTION_IDLE" env-default:"0s"`          // Пример: "30m". Если клиент неактивен это время, сервер посылает GOAWAY. 0s = бесконечность (gRPC default).
	KeepAliveMaxConnectionAge      time.Duration `yaml:"keepAliveMaxConnectionAge" env:"GRPC_KEEP_ALIVE_MAX_CONNECTION_AGE" env-default:"0s"`            // Пример: "2h". Максимальное время жизни соединения. 0s = бесконечность (gRPC default).
	KeepAliveMaxConnectionAgeGrace time.Duration `yaml:"keepAliveMaxConnectionAgeGrace" env:"GRPC_KEEP_ALIVE_MAX_CONNECTION_AGE_GRACE" env-default:"0s"` // Пример: "30s". Дополнительное время после MaxConnectionAge. 0s = бесконечность (gRPC default).
	KeepAliveServerTime            time.Duration `yaml:"keepAliveServerTime" env:"GRPC_KEEP_ALIVE_SERVER_TIME" env-default:"2h"`                         // Если клиент неактивен это время, сервер пингует клиента. (gRPC default: 2h)
	KeepAliveServerTimeout         time.Duration `yaml:"keepAliveServerTimeout" env:"GRPC_KEEP_ALIVE_SERVER_TIMEOUT" env-default:"20s"`                  // Время ожидания ответа на пинг перед закрытием соединения. (gRPC default: 20s)

	// Политика принудительного KeepAlive (grpc.keepalive.EnforcementPolicy - KAEP)
	// Защищает сервер от злоупотреблений настройками keep-alive со стороны клиента.
	KeepAliveEnforcementPolicyMinTime             time.Duration `yaml:"keepAliveEnforcementPolicyMinTime" env:"GRPC_KEEP_ALIVE_ENFORCEMENT_MIN_TIME" env-default:"5m"`                             // Минимальный интервал между пингами от клиента. (gRPC default: 5m)
	KeepAliveEnforcementPolicyPermitWithoutStream bool          `yaml:"keepAliveEnforcementPolicyPermitWithoutStream" env:"GRPC_KEEP_ALIVE_ENFORCEMENT_PERMIT_WITHOUT_STREAM" env-default:"false"` // Разрешать ли пинги без активных потоков. (gRPC default: false)

	// Лимиты размеров сообщений
	MaxReceiveMessageSize int `yaml:"maxReceiveMessageSize" env:"GRPC_MAX_RECEIVE_MESSAGE_SIZE" env-default:"4194304"` // 4MB (gRPC default: 4MB)
	MaxSendMessageSize    int `yaml:"maxSendMessageSize" env:"GRPC_MAX_SEND_MESSAGE_SIZE" env-default:"0"`             // 0 для использования gRPC default (math.MaxInt32)

	// Лимиты потоков и параллелизма
	MaxConcurrentStreams  uint32 `yaml:"maxConcurrentStreams" env:"GRPC_MAX_CONCURRENT_STREAMS" env-default:"0"`    // Максимальное количество одновременных потоков на одном HTTP/2 соединении. 0 для gRPC default (math.MaxUint32).
	InitialWindowSize     int32  `yaml:"initialWindowSize" env:"GRPC_INITIAL_WINDOW_SIZE" env-default:"0"`          // Начальный размер окна для потока. 0 для gRPC default (64KB).
	InitialConnWindowSize int32  `yaml:"initialConnWindowSize" env:"GRPC_INITIAL_CONN_WINDOW_SIZE" env-default:"0"` // Начальный размер окна для соединения. 0 для gRPC default.

	// Размеры буферов (для настройки производительности)
	ReadBufferSize  int `yaml:"readBufferSize" env:"GRPC_READ_BUFFER_SIZE" env-default:"32768"`   // 32KB (gRPC default)
	WriteBufferSize int `yaml:"writeBufferSize" env:"GRPC_WRITE_BUFFER_SIZE" env-default:"32768"` // 32KB (gRPC default)

	// Регистрация стандартных сервисов
	EnableHealthCheckService bool `yaml:"enableHealthCheckService" env:"GRPC_ENABLE_HEALTH_CHECK_SERVICE" env-default:"true"` // Автоматически регистрировать стандартный Health Check сервис.
	EnableReflectionService  bool `yaml:"enableReflectionService" env:"GRPC_ENABLE_REFLECTION_SERVICE" env-default:"true"`    // Автоматически регистрировать Reflection сервис (для grpcurl, etc.).

	// "Вежливое" завершение работы (Graceful Shutdown)
	GracefulShutdownTimeout time.Duration `yaml:"gracefulShutdownTimeout" env:"GRPC_GRACEFUL_SHUTDOWN_TIMEOUT" env-default:"30s"` // Таймаут для ожидания завершения активных RPC перед принудительной остановкой.

	// Таймаут на установку соединения (до TLS handshake и создания потоков)
	ConnectionTimeout time.Duration `yaml:"connectionTimeout" env:"GRPC_CONNECTION_TIMEOUT" env-default:"120s"` // (gRPC default: 120s)

}

type GrpcClient struct {
	Host      string `yaml:"host" env:"HOST,required"`
	Port      int    `yaml:"port" env:"PORT,required"`
	UserAgent string `yaml:"userAgent" env:"USER_AGENT"`

	// TLS
	EnableTLS          bool   `yaml:"enableTLS" env:"ENABLE_TLS" env-default:"false"`
	CACertFile         string `yaml:"caCertFile" env:"CA_CERT_FILE"`                 // Путь к CA сертификату сервера
	ClientCertFile     string `yaml:"clientCertFile" env:"CLIENT_CERT_FILE"`         // Путь к клиентскому сертификату (для mTLS)
	ClientKeyFile      string `yaml:"clientKeyFile" env:"CLIENT_KEY_FILE"`           // Путь к приватному ключу клиента (для mTLS)
	ServerNameOverride string `yaml:"serverNameOverride" env:"SERVER_NAME_OVERRIDE"` // Переопределение имени сервера в TLS сертификате

	// Параметры повторных попыток подключения (для метода Connect)
	ConnectMaxAttempts       int           `yaml:"connectMaxAttempts" env:"CONNECT_MAX_ATTEMPTS" env-default:"5"`
	ConnectInitialBackoff    time.Duration `yaml:"connectInitialBackoff" env:"CONNECT_INITIAL_BACKOFF" env-default:"250ms"`
	ConnectMaxBackoff        time.Duration `yaml:"connectMaxBackoff" env:"CONNECT_MAX_BACKOFF" env-default:"5s"`
	ConnectBackoffMultiplier float64       `yaml:"connectBackoffMultiplier" env:"CONNECT_BACKOFF_MULTIPLIER" env-default:"2.0"`

	// Тайм-аут на установку соединения (для каждой отдельной попытки grpc.DialContext)
	DialTimeout time.Duration `yaml:"dialTimeout" env:"DIAL_TIMEOUT" env-default:"5s"`

	// KeepAlive параметры
	KeepAliveTime       time.Duration `yaml:"keepAliveTime" env:"KEEP_ALIVE_TIME" env-default:"30s"`
	KeepAliveTimeout    time.Duration `yaml:"keepAliveTimeout" env:"KEEP_ALIVE_TIMEOUT" env-default:"20s"`
	PermitWithoutStream bool          `yaml:"permitWithoutStream" env:"PERMIT_WITHOUT_STREAM" env-default:"true"`

	// Размеры сообщений
	MaxRecvMsgSize int `yaml:"maxRecvMsgSize" env:"MAX_RECV_MSG_SIZE" env-default:"4194304"` // 4MB
	MaxSendMsgSize int `yaml:"maxSendMsgSize" env:"MAX_SEND_MSG_SIZE" env-default:"4194304"` // 4MB
}

type Grpc struct {
	Server  GrpcServer   `yaml:"server"`
	Clients []GrpcClient `yaml:"clients"`
}

func MustLoad[TConfig Config]() (*TConfig, *Env) {
	path := fetchConfigPath()

	if path == "" {
		panic("Путь конфига не найден")
	}

	if _, err := os.Stat(path); err != nil {
		panic("Файл конфига не найден")
	}

	var cfg TConfig

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("Ошибка загрузки конфига: " + err.Error())
	}

	env, err := NewEnv(cfg.Env())
	if err != nil {
		panic("Ошибка создания окружения: " + err.Error())
	}

	return &cfg, env
}

func fetchConfigPath() string {
	var result string
	flag.StringVar(&result, "config", "", "Путь до файла конфига")
	flag.Parse()

	if result == "" {
		result = os.Getenv("CONFIG_PATH")
	}

	return result
}
