# Luzifer / envrun

`envrun` is a small helper utility I wrote for myself to debug programs and scripts during their development expecting environment variables to be set to special values. Sure there is [gin](https://github.com/codegangsta/gin) for go webservers doing the same but I wanted something also for commandline utilities.

It reads a `.env` file (configurable) from the current directory and then either takes its own environment variables or a clean set and adds the env variables found in `.env` to it. The resulting set is passed to the command you put as arguments to `envrun`.

## Examples

To visualize the effect of the utility the test command is `python test.py` with this simple python script:

```python
import os

for k in os.environ.keys():
  print "{} = {}".format(k, os.environ[k])
```

It just prints the current environment to `STDOUT` and exits.

```bash
# cat .env
MY_TEST_VAR=hello world
ANOTHER_VAR=foo

# python test.py | grep MY_TEST_VAR
## No output on this command

# envrun --help
Usage of envrun:
      --clean[=false]: Do not pass current environment to child process
      --env-file=".env": Location of the environment file
      --q[=false]: Suppress informational messages from envrun

# envrun python test.py | grep MY_TEST_VAR
MY_TEST_VAR = hello world

# envrun python test.py | wc -l
      45

# envrun --clean python test.py | wc -l
       3

# envrun --clean python test.py
__CF_USER_TEXT_ENCODING = 0x1F5:0x0:0x0
ANOTHER_VAR = foo
MY_TEST_VAR = hello world
```

## Encrypted `.env`-file

In case you don't want to put the environment variables into a plain text file onto your disk you can use an AES256 encrypted file and provide a password to `envrun`:

```bash
# echo 'MYVAR=myvalue' | openssl enc -e -aes-256-cbc -pass pass:justatest -base64 -out .env

# cat .env
U2FsdGVkX18xcVIMejjwWzh1DppzptJCHhORH/JDj10=

# envrun -p justatest --clean -- env
MYVAR=myvalue
2017/03/21 16:34:57 Process exitted with code 0
```
