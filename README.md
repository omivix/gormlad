### Example

```golang
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	Logger: gormlad.New(lad.L()).LogMode(logger.Info),
})
```
