#! /usr/bin/env python3
# -*- coding: utf-8 -*-

import sys

# Ce programme sert à modifier le fichier app.go et modifier les path.
# Si je suis en prod ou en dev

def check_args(*args, **kwargs):
  if len(sys.argv) != 2:
    help()
    sys.exit(0)

  if sys.argv[1] == "prod":
    return "prod"
  elif sys.argv[1] == "dev":
    return "dev"
  else:
    help()
    sys.exit(0)

def prod_change(*args, **kwargs):
  slurp_splitted = kwargs.get("slurp_splitted")

  for num, line in enumerate(slurp_splitted):
      if "confPath := " in line:
        slurp_splitted[num] = '\tconfPath := "/etc/gwmpd"'
      elif "confFilename := " in line:
        slurp_splitted[num] = '\tconfFilename := "gwmpd"'
      elif "logFilename := " in line:
        slurp_splitted[num] = '\tlogFilename := "/var/log/gwmpd/error.log"'

  return slurp_splitted

def dev_change(*args, **kwargs):
  slurp_splitted = kwargs.get("slurp_splitted")

  for num, line in enumerate(slurp_splitted):
      if "confPath := " in line:
        slurp_splitted[num] = '\tconfPath := "cfg/"'
      elif "confFilename := " in line:
        slurp_splitted[num] = '\tconfFilename := "gwmpd_sample"'
      elif "logFilename := " in line:
        slurp_splitted[num] = '\tlogFilename := "error.log"'

  return slurp_splitted

def edit_file(*args, **kwargs):
  comp = kwargs.get("compile")
  slurp_splitted = kwargs.get("slurp_splitted")

  if comp == "prod":
    return prod_change(slurp_splitted=slurp_splitted)
  elif comp == "dev":
    return dev_change(slurp_splitted=slurp_splitted)
  else:
    # Ooooooopppssss
    print("Something was wrong. Unknow comp: {}".format(comp))
    sys.exit(1)

def read_file(*args, **kwargs):
  filename = kwargs.get("filename")
  slurp = ""

  with open(filename, 'r') as fn:
    slurp = fn.read()

  return str.split(slurp, "\n")

def write_file(*args, **kwargs):
  filename = kwargs.get("filename")
  slurp_splitted = kwargs.get("slurp_splitted")

  slurp = "\n".join(slurp_splitted)
  with open(filename, 'w') as fn:
    fn.write(slurp)

def help():
  print("{} dev|prod".format(sys.argv[0]))
  print("  dev\tChange app.go to use dev env")
  print("  prod\tChange app.go to use prod env")

def main(*args, **kwargs):
  fn = "app.go"

  comp = check_args()
  slurp_splitted = read_file(filename=fn)
  slurp_splitted = edit_file(compile=comp, slurp_splitted=slurp_splitted)
  write_file(filename=fn, slurp_splitted=slurp_splitted)

if __name__ == "__main__":
  main()
