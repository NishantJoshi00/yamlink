#!/bin/zsh

# first args is the path of the `wayshell` executable
# second args is the path of the `wayshell` config file

__ERROR=false

# verify if the first argument is a file
if [[ ! -f $1 ]]; then
  echo "The first argument must be a file"
  __ERROR=true
fi

# verify if the second argument is a file
if [[ ! -f $2 ]]; then
  echo "The second argument must be a file"
  __ERROR=true
fi


if [[ $__ERROR == false ]]; then
  EXE=$1
  CONFIG=$2
  
  wayshell-expr () {
    if [[ "$BUFFER" == "s/"* ]]; then
  
      # here the query is `/.*` e.g ( if the line looks like `s/hello/world` then the query is `/hello/world` )
      query=$(echo $BUFFER | sed 's/^s\(\/\S*\).*/\1/')
      args=$(echo $BUFFER | sed 's/^s\/\S*\(.*\)/\1/')
      BUFFER="$(CONFIG_FILE=$CONFIG $EXE $query)$args"
    fi
  
    zle .accept-line
  }
  
  zle -N accept-line wayshell-expr
fi


