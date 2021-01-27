# script-wrapper
Wrapper for script file

command:
```shell script
# encrypt script
./swrapper -m enc -p plain.sh -c plain.sh.enc

# execute script
./swrapper -m run -c plain.sh.enc

# args: 
-p: Plain text file (Default: plain.sh)
-c: Cipher text file (Default: data.dat)
-m: Execute mode[enc/run] {enc: 'encrypt script', run: 'run script'}
```


