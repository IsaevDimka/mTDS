__MetaTDS v.2 CHANGE LOG__

> __04.12.2018__
+ Core functions and logic
+ GIN-GO Framework in use
+ Redirect functions (old format)

> __05.12.2018__
+ Config `settings.ini` file reader / writer
+ Partial ports from LUA language TDS v.1
+ Redirect functions (new format)

> __06.12.2018__
+ Self format for RedisDB
+ Self format for respond
+ Debug utilities

> __07.12.2018__
+ Flow add/update API via HTTP
+ Click info API + formats for back-requests from prelands
+ Click generation
+ Types and methods for models

> __09.12.2018__
+ GIN-GO Framework changed to ECHO Minimalist
+ Functional dropped to file modules
+ Optimizing and refactoring some parts
+ test build for Ubuntu 18.04

> __10.12.2018__
+ Support for backcalls from Prelands
+ Extended info by all Clicks

> __11.12.2018__
+ JSON format support for all features

> __12.12.2018__
+ JSON format requests / responds
+ DEP package manager support
+ test build via `gccgo` compiler for Ubuntu 18.04

> __14.12.2018__
+ Github.com working repo + branches
+ Additional features for HIT export
+ Functional dropped to packages

> __16.12.2018__
+ Telegram notifications for statistic / usage
+ Tested with GCCGO compiler over Ubuntu 18.04
+ Refactored directory structure
+ Minor fixes

> __17.12.2018__
+ Time bench and statistics
+ Dynamic reload .ini file to get settings online
+ Sending custom stats to Telegram Bot
+ Import flows into Redis

> __18.12.2018__
+ Failover refusing connections withing 5sec. IN/OUT
+ Redis crash and recover support
+ Import Flows to Redis
+ Save data to tests-log file `./tdstest.log`
+ Advanced setting cookie
+ Time usage statistics

> __19.12.2018__
+ fx 30 seconds drop connection
+ Recovering when Cookie reading
+ Redis temp func removeing objects 1h/rate
+ Export list of clicks 
+ Creating click with path to append
  ClickHouse export format
  

> __23.12.2018__
+ Improvements for Telegram Bot
+ Saving clicks from RedisDB to a files dir ./clicks/XXXXXXXXXX.json
+ Improvements for Statistics and memory usage params
+ Garbage collector implements
+ Mass testing analasys and fixes bugs depends on

> __24.12.2018__
+ for more stabiliy and graceful deployment
  `settings.ini` is not in a package anymore, now after deployment
  you should rename `settings.dev` or `settings.prod` to `settings.ini
+ fixed bug with nil pointer when telegram can't connect to proxy
+ added library for memory/system monitoring  

> __25.12.2018__
+ testing issues (too many open files)
+ count open files within statistic
+ absolute file path for `settings.ini`
+ added service configuration for system.d (thanks to `@pfirsov`) 

> __26.12.2018__
+ Sending to API
+ Re-sending in case of failure to API
+ Log files optional to STDOUT 
+ Improvements for `system.d` start script
+ Statistics improvements + average time for respond

> __27.12.2018__
+ Recieving flows from API for startup and for work
+ Re-sending in case of failure to API (clicks) tested
+ Additional time stamp saved to `last.update.time` 
+ Improvements for `system.d` start script
+ Statistics improvements + average time for respond
+ Options for flows and clicks @`settings.ini`

> __30.12.2018__
+ Improved performance for responsing to context
+ Additional statistics /stat, /conf
+ Small refactoring
+ Additional configurable settings @`settings.ini`
+ Production testing

> __31.12.2018 22:29 (коммит новогодний)__
+ Improved average time calculatin
+ Improved saving click with redis HMSET
+ Improved getting infos from all triiggers
  lands / prelands / flows / clicks
+ reused profiling information
+ Humanize
+ fx with count of flows

> __01.01.2019__
+ tests to increase performance of randomization
+ additional configuration & management tools
+ memory profiling
+ cpu profiling
+ string generator fix
+ resend statistic send and file saving refactoring

> __02.01.2019__

+ Multi deleting from Redis by exec One command
+ File I/O reading improvements@go routines

> __03.01.2019__

+ Multi deleting from Redis by exec One command
+ File I/O reading improvements@go routines
+ final testing on high-load with concurency 10000rps

> __06.01.2019__

+ Multi deleting from Redis improvements
+ Reconstruced URL helper
+ final testing on high-load with concurency 1000rps/1billion requests

> __07.01.2018__
+ Improved performance for importing flows
+ Additional statistics /stat, /conf
+ Small refactoring for init.go
+ Additional configurable settings @`settings.ini`
+ Production testing on +10000 rps

__FEATURED TODO__
+ Implementing pixel integration 
+ Mass load testing 
