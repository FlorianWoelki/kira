This document contains possible exploits of kira and how they are fixed. All possible
exploits can be found in the `/exploits` directory which includes python files. Furthermore,
kira is running in a Dockerfile which has a lot of potential security fixes out of the box.

If you have any other potential security flaw, make sure to create an issue (and maybe
possible fix with a pull request) with the description of the exploit.

# Endless Loop

The endless loop exploit simply describes a never ending loop and can lead to potential
memory exhaustion or a job will never quit or end.

Example Code:
```py
while True:
  print("Hello World")

```

This exploit has been fixed through limiting the amount of memory for each process and
adding a timeout that is not dependent on the API or running job.

# Fork

An endless loop that constantly forks the OS, which can lead to memory exhaustion.

Example Code:
```py
import os

while True:
  try:
    os.fork()
  except:
    pass
```

This exploit has been fixed through setting the maximum processes to 256 by default. This
fix also fixes exploits like `:(){ :|: &}:;`.

# Exhausting Space

Exhausting space through a file based attack can be really crucial and can lead to several
issues for kira. The following code creates a massive file which might exhaust the memory
or disk space.

Example Code:
```py
with open("test", "w") as f:
  n = 2**24
  f.write("This is a test\n" * n)
```

The fix for this exploit has been achieved through setting the maximum files to 2048 by
default. Furthermore, all the temporary created directories will be deleted after each
execution which resists out of drive space attacks.

# External Resources

Downloading specific files from other servers can be malicious. This can harm kira in a
crucial way.

Example Code:
```py
import urllib.request
with urllib.request.urlopen("https://somecode.com") as response:
  contents = response.read()
  print(contents)
```

That's why the fix for this exploit was to disable the network interface completely.

# Writing to other files

A user that uses this service could always write to other main files that were not part of their code directory. This was a crucial bug, because malicious code could have been added to other files.

Example Code:
```py
f = open("./languages/python/run.sh", "w")
print(f.write('Test'))
```

This has been fixed by logging into the current user's shell and executing the files from there.
