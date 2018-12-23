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
+ Telegram adapter
+ Time bench for Each external calls/handlers
+ Statistic export to Telegram MetaTDSBot
+ Changed project directory tree

> __19.12.2018__
+ fx 30 seconds drop connection
+ Recovering when Cookie reading
+ Redis temp func removeing objects 1h/rate
+ Export list of clicks 
+ Creating click with path to append
  ClickHouse export format

> __18.12.2018__
+ Failover refusing connections withing 5sec. IN/OUT
+ Redis crash and recover support
+ Import Flows to Redis
+ Save data to tests-log file `./tdstest.log`
+ Advanced setting cookie
+ Time usage statistics

> __17.12.2018__
+ Time bench and statistics
+ Dynamic reload .ini file to get settings online
+ Sending custom stats to Telegram Bot
+ Import flows into Redis

> __16.12.2018__
+ Telegram notifications for statistic / usage
+ Tested with GCCGO compiler over Ubuntu 18.04
+ Refactored directory structure
+ Minor fixes

__FEATURED TODO__
+ Tool to import flows (mass/single)
+ Implementing pixel integration 
+ Sending stats to MetaData API
+ Saving to files is ready and depens on previous article
+ Mass load testing 
+ Check dir for existence when saving clicks to ./clicks/XXXXXXX.json
