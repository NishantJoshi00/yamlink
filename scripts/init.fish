#!/bin/fish


if test (count $argv) -eq 2 -a -f $argv[1] -a -f $argv[2]
  echo "Initializing fish shell with exec file: $argv[1] and config file: $argv[2]"

  set exec_file $argv[1]
  set config_file $argv[2]

  function modify_command
    set argv (commandline)
    # check if argv has s/ at the beginning
    if test (string match -r '^s/' $argv[1])
      set query (string replace -r '^s(\S*).*' '$1' $argv[1])
      set args (string replace -r '^s\S*(.*)' '$1' $argv[1])
      set output (CONFIG_FILE=$config_file $exec_file $query)
      commandline -r (string join '' $output $args)
    end

    commandline -f execute
  end

  bind \r modify_command

else
  echo "Failed while initializing fish shell. Please provide two files as arguments."
end

