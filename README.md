__FEATURED TODO__
+ Development testing
+ Production testing

__LAST CHANGES__

18.12.2018
+ Failover refusing connections withing 5sec. IN/OUT
+ Redis crash and recover support
+ Import Flows to Redis
+ Save data to tests-log file `./tdstest.log`
+ Advanced setting cookie
+ Time usage statistics

17.12.2018
+ Time bench and statistics
+ Dynamic reload .ini file to get settings online
+ Sending custom stats to Telegram Bot
+ Import flows into Redis

16.12.2018
+ Telegram notifications for statistic / usage
+ Tested with GCCGO compiler over Ubuntu 18.04
+ Refactored directory structure
+ Minor fixes

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
