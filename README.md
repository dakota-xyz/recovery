# Recovery tool

This recovery tool can be used by Dakota customers to recover their keys.
It requires a **decrypted** backup shard, the client shard, and a JSON file with the key mappings.

The format of the key mappings JSON file should be similar to

```json
{
  "organization_id": "e65ccfaa-01b7-4a00-aec5-0fcc25d7eba7",
  "keys": [
    {
      "address_sub_id": "21d3969c-8a56-46d9-be38-b53c21294e54",
      "network_id": "solana-mainnet",
      "curve": "ELLIPTIC_CURVE_ED25519"
    },
    {
      "address_sub_id": "a09d020a-1bc8-47b6-a208-0fd47cd05e66",
      "network_id": "ethereum-mainnet",
      "curve": "ELLIPTIC_CURVE_SECP256K1"
    }
  ]
}

```

# Build

## Install Go 1.21 or higher:

Currently, recovery uses Go 1.21 to compile the code.

Install [Go 1.21](https://go.dev/doc/install) by following instructions there, e.g.:

```sh
wget https://go.dev/dl/go1.21.3.linux-amd64.tar.gz
sudo tar xzvf go1.21.3.linux-amd64.tar.gz -C /usr/local/
export PATH=$PATH:/usr/local/go/bin
```

Verify the installation by typing `go version` in your terminal.

```sh
$ go version
go version go1.21.3 darwin/amd64
```

## Build Recovery

In order to build `recovery` you need the source code. Either [download the source of a release](https://github.com/dakota-xyz/recovery/releases) or [clone the git repository](https://github.com/dakota-xyz/recovery).

Build `recovery` from the source code:

``` 
cd recovery
make build
```

After building, you should see a new executable file `recovery/build/recovery`.

# Example usage

```bash
$ ./recovery -h
Usage of ./recovery:
  -keymap string
        Location of the JSON file containing the key map
  -shard1 string
        Location of the file containing the first shard
  -shard2 string
        Location of the file containing the second shard
  -target string
        Target CSV file (default "keys.csv")
$ ./recovery -shard1 shard1.bin -shard2 shard2.bin -keymap backup.json
2023/11/09 03:20:10 INFO Initiating recovery
2023/11/09 03:20:10 INFO Recovery complete. Results saved to keys.csv
$ cat keys.csv
Network,Address,PrivateKey
solana-mainnet,4tZpnxbJbkCDFFCpbmb4y7wsH366kxeb57R8owi67qi8,2tFuN9PCkTYsDV6rq8RauJZEmyBs7x8rLoSAFYD5JcQMCzVzStq45VeUVDDghGqXaYm8muC8YECzgoqTkyPph8gp
ethereum-mainnet,0x0D7ad5799E3DB77c8258b9700E4f94Fcb092C64B,0x222d55b028c7896058d28af1d44c55d45264c470f2a93e7b013076e68b7bfa25
```
