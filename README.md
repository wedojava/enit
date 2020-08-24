#Enit

Encrypt it!  
Encrypt any string in pattern of [key, value], and have a key argument for encoding is blank by default, you can set it via `-k <your encoding_key>`, anyway, it encrypt [key, value] via aes.  

# Usage
Run `./enit` to read help, the fast begin case is below:
```
$ ./enit set <yourkey> <yourvalue>
$ ./enit get <yourkey>
// Output:
$ <yourkey> = <yourvalue>
```
Also, you can add key for encoding secret file.  
```
$ ./enit set -k <encoding_key> <yourkey> <yourvalue>
$ ./enit get -k <encoding_key> <yourkey>
// Output:
$ <yourkey> = <yourvalue>
```

# Details
The config file is located at `<your_home_dir>/.enit` by default.  

More feature to be continue...  

Enjoy it!  

EOF
