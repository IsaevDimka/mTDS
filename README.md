__FEATURED TODO__
+ Development testing
+ Production testing

__LAST CHANGES__

> __31.12.2018__
 + Improved average time calculatin
 + Improved saving click with redis HMSET
 + Improved getting infos from all triiggers
   lands / prelands / flows / clicks
 + reused profiling information

> __30.12.2018__
+ Improved performance for responsing to context
+ Additional statistics /stat, /conf
+ Small refactoring
+ Additional configurable settings @`settings.ini`
+ Production testing

> __27.12.2018__
+ Recieving flows from API for startup and for work
+ Re-sending in case of failure to API (clicks) tested
+ Additional time stamp saved to `last.update.time` 
+ Improvements for `system.d` start script
+ Statistics improvements + average time for respond
+ Options for flows and clicks @`settings.ini`

> __26.12.2018__
+ Sending to API
+ Re-sending in case of failure to API
+ Log files optional to STDOUT 
+ Improvements for `system.d` start script
+ Statistics improvements + average time for respond

> __25.12.2018__
+ testing issues (too many open files)
+ count open files within statistic
+ absolute file path for `settings.ini`
+ added service configuration for system.d (thanks to `@pfirsov`) 

> __24.12.2018__
+ for more stabiliy and graceful deployment
  `settings.ini` is not in a package anymore, now after deployment
  you should rename `settings.dev` or `settings.prod` to `settings.ini
+ fixed bug with nil pointer when telegram can't connect to proxy
+ added library for memory/system monitoring  

> __23.12.2018__
+ Improvements for Telegram Bot
+ Saving clicks from RedisDB to a files dir ./clicks/XXXXXXXXXX.json
+ Improvements for Statistics and memory usage params
+ Garbage collector implements
+ Mass testing analasys and fixes bugs depends on

__INSTALLATION__

1. clone this repository like this:
> `$ git clone https://github.com/CpaMonstr/metatds.git`

2. in folder `metacds` run following comand
> `$ dep init`

if you have no `dep` installed, just use
> `$ curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh`

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
