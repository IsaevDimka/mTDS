__INSTALLATION__

1. clone this repository like this:
> `$ git clone https://github.com/CpaMonstr/metatds.git`

2. in folder `metacds` run following comand
> `$ dep init`

3. edit configuration file `settings.ini` to make sure
your settings are correct

4. Run the following command in project folder

for Linux
> `$ ./make.sh` 

or for Windows
> `$ make.cmd`

it will build project and automaticly run it.

Additional information you can find in /docs folder,
there are format specifications and usage examples

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

__TODO FEATURES__

+ utm_source convertion to subs
+ ask in redis for uniquity of generated CLICK_HASH
+ testing for production