# Summer IOC Framework

summer是一个开箱即用的ioc框架。

功能：
1. 根据类型或者tag解决组件之间的依赖
2. 提供toml插件，根据tag直接写入组件的属性，不需要设计config的struct。
3. 提供生命周期的管理，生命周期分为 
    1. init 执行一些不需要依赖其他组件的初始化动作,一般是初始化组件基本信息，如链接数据库，链接zookeeper 
    2. ready 组件执行ready表示本身组件已经可被调用,所有依赖组件的功能也可被调用。
    此时可以执行一些需要依赖组件的初始化动作 
    3. destroy 此时将会销毁或者卸载组件
4. 定义了插件调用的时间锚点，可以自定义插件.

具体使用示例：
1. 基本依赖与生命周期
   - 场景:每隔一段时间从redis同步数据到mysql的小程序
   - 代码 (具体见 *example/sync-data* 目录): 
      
       1. RedisProvider 负责提供redis链接
        ```go
           func init() {
               summer.Put(&RedisProvider{})
           }
           type RedisProvider struct {
               Client *redis.Client
           }
           
           func (provider *RedisProvider) Init() {
               provider.Client = redis.NewClient(&redis.Options{
                   Addr: Conf.RedisAddr,
               })
               err := provider.Client.Ping().Err()
               if err != nil {
                   panic(err)
               }
           }
           
           func (provider *RedisProvider) Provide() (client *redis.Client) {
               return provider.Client
           }
        ```
       2. DatabaseProvider 负责提供数据库链接
        ```go
        func init() {
            summer.Put(&DatabaseProvider{})
        }
        type DatabaseProvider struct {
            DB *sql.DB
        }
        
        func (provider *DatabaseProvider) Init() {
            conn, err := sql.Open("mysql", Conf.MysqlDSN)
            if err != nil {
                panic(err)
            }
            provider.DB = conn
        }
        
        func (provider *DatabaseProvider) Provide() (db *sql.DB) {
            return provider.DB
        }
        ```
       3. SyncDataWorker 定时执行同步数据操作
       ```go
        const key = "/key"
        
        func init() {
            summer.Put(&SyncDataWorker{})
        }
        
        type SyncDataWorker struct {
            RedisProvider    *RedisProvider `sm:"*"`
            DatabaseProvider *DatabaseProvider `sm:"*"`
            redisClient      *redis.Client
            db               *sql.DB
        }
        
        func (worker *SyncDataWorker) Run() {
            for {
                if result := worker.redisClient.Get(key); result.Err() == nil {
                    worker.db.Exec("update `test_table` set `text` = ? where `key` = ? ", result.String(), key)
                } else {
                    log.Println(result.Err())
                }
                time.Sleep(time.Minute)
            }
        }
        func (worker *SyncDataWorker) Ready() {
            worker.redisClient = worker.RedisProvider.Provide()
            worker.db = worker.DatabaseProvider.Provide()
            go worker.Run()
        }

        ```
  
   