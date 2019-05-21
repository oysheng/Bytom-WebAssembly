# Bytom-WebAssembly
It is a project for Bytom WebAssembly

## Prepare
```sh
git clone https://github.com/oysheng/Bytom-WebAssembly.git $GOPATH/src/github.com/bytom-community/wasm
```

## Build

Need Go version 1.11

```sh
cd $GOPATH/src/github.com/bytom-community/wasm
#default build
GOOS=js GOARCH=wasm go build -o main.wasm
#mini build
GOOS=js GOARCH=wasm go build -tags=mini -o main.wasm 
```


## WebAssembly JS Function
##### mini build
>createKey\
resetKeyPassword \
signTransaction

##### default build
>createKey \
resetKeyPassword \
createAccount \
createAccountReceiver \
signTransaction \
signMessage \
convertArgument \
createPubkey
