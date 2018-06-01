[![Go Report Card](https://goreportcard.com/badge/github.com/Luzifer/envrun)](https://goreportcard.com/report/github.com/Luzifer/envrun)
![](https://badges.fyi/github/license/Luzifer/envrun)
![](https://badges.fyi/github/downloads/Luzifer/envrun)
![](https://badges.fyi/github/latest-release/Luzifer/envrun)

# Luzifer / envrun

`envrun` is a small helper utility to inject environment variables stored in a file into processes.

It reads a `.env` file (configurable) from the current directory and then either takes its own environment variables or a clean set and adds the env variables found in `.env` to it. The resulting set is passed to the command you put as arguments to `envrun`.

## Examples

To visualize the effect of the utility the test command is `python test.py` with this simple python script:

```python
import os

for k in os.environ.keys():
  print "{} = {}".format(k, os.environ[k])
```

It just prints the current environment to `STDOUT` and exits.

```console
$ cat .env
MY_TEST_VAR=hello world
ANOTHER_VAR=foo

$ python test.py | grep MY_TEST_VAR
## No output on this command

$ envrun --help
Usage of envrun:
      --clean                  Do not pass current environment to child process
      --encryption string      Encryption method used for encrypted env-file (Available: openssl-md5) (default "openssl-md5")
      --env-file string        Location of the environment file (default ".env")
      --log-level string       Log level (debug, info, warn, error, fatal) (default "info")
  -p, --password string        Password to decrypt environment file
      --password-file string   Read encryption key from file
      --q                      Suppress informational messages from envrun (DEPRECATED, use --log-level=warn)
      --version                Prints current version and exits

$ envrun python test.py | grep MY_TEST_VAR
MY_TEST_VAR = hello world

$ envrun python test.py | wc -l
      45

$ envrun --clean python test.py | wc -l
       3

$ envrun --clean python test.py
__CF_USER_TEXT_ENCODING = 0x1F5:0x0:0x0
ANOTHER_VAR = foo
MY_TEST_VAR = hello world
```

## Encrypted `.env`-file

In case you don't want to put the environment variables into a plain text file onto your disk you can use an encrypted file and provide a password to `envrun`:

### GnuPG symmetric encryption

In this example an armored (`-a`) encryption is used. This is not required and you can leave out the `-a` flag.

```console
$ echo "MYVAR=myvalue" | gpg --passphrase justatest --batch --quiet --yes -c -a -o .env

$ cat .env
-----BEGIN PGP MESSAGE-----

jA0ECQMCIsGVKNlJw1Py0kMB542XJvekKyuPi2LHQrnFlhD5ALm6orvE3WFAzp7D
kAisTMr10fmjLuENfQhxqd9MB0Kd2mfd3b1mgOzei5IMDLJc
=7k9M
-----END PGP MESSAGE-----

$ envrun -p justatest --encryption gpg-symmetric --clean -- env
MYVAR=myvalue
INFO[0000] Process exitted with code 0
```

### OpenSSL AES256 encryption

Pay attention on the `-md md5` flag: OpenSSL 1.1.0f and newer uses an incompatible hasing algorithm for the passwords!

```console
$ echo 'MYVAR=myvalue' | openssl enc -e -aes-256-cbc -pass pass:justatest -md md5 -base64 -out .env

$ cat .env
U2FsdGVkX18xcVIMejjwWzh1DppzptJCHhORH/JDj10=

$ envrun -p justatest --clean -- env
MYVAR=myvalue
2017/03/21 16:34:57 Process exitted with code 0
```
