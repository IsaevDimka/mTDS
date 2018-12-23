__FEATURED TODO__
+ Development testing
+ Production testing

__LAST CHANGES__

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
