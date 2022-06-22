function reverse_cmd_history_loopup
  #history | peco --layout bottom-up | xargs -I _ fish -c '_'
  commandline -r (history | peco --layout bottom-up)
end


function search_current_dir
  read -P "Search Query: " FOO
  if [ $FOO ]
    set search_result (rg --count-matches --auto-hybrid-regex --trim $FOO | peco)
    set search_result_arr (string split ':' $search_result)
    commandline -i (string escape $search_result_arr[1])
  end
  commandline -i ''
end


function login
  echo "logging into tools"
end

function inject_previous_command
  commandline -i (string escape $previous_command)
end

# bind \cR reverse_cmd_history_loopup
# bind \cS search_current_dir
# bind \cF search_current_dir
# bind \cW login


# bind \e\[1\;9A inject_previous_command