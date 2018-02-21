Setup to run examples
-----------------------------------

First get "geth" up and running on your system.

Second, create/modify a JSON file, ./setup.json, with your configuration data

```json
{
	"keyFile": "../UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
	"keyFilePassword": "password",
	"gethServer": "ws://192.168.0.157:8546",
	"contractAddress": "0xe705ab560794cf4912960e5069d23ad6420acde7"
}
```

in this directory.  `cd` into the example.  Compile and run.  The examples will open ../setup.json.
The default file is configured for my use in the testnet.

Note
------------

Your geth will probably be at 127.0.0.1, not my class 'C' network address of 192.168.*.

Don't try to run on an IP version 6 address - that will not work.  On some Mac's the 
defauilt for "localhost" is an IP version 6 address - so you may need to use 127.0.0.1
to force it to be a IP verison 4 address.

Run the event watcher before running the event generator.

