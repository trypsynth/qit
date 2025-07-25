def usage : Nil
  puts <<-USAGE
  Qit - Quin's tiny Git helper.
  Usage: qit <command> [<args>...]
  Available commands:
    acp <message>: add all files, commit with the specified message, and push.
    amend <message>: amend the last commit with a new message.
    last: view the most recent commit message.
    log: show the commit log in a readable format.
    undo: undo the last commit while keeping changes intact.
  USAGE
end

def git(*args : String) : Nil
  abort "Command 'git #{args.join(" ")}' failed." unless Process.run("git", args, output: STDOUT, error: STDERR).success?
end

def get_commit_message(prefix : String = "") : String
  if ARGV.size < 2
    desc = prefix.empty? ? "commit message" : "#{prefix} commit message"
    STDERR.puts "Missing#{desc}."
    usage
    exit 1
  end
  ARGV[1..].join(" ")
end

if ARGV.empty?
  usage
  exit 1
end
command = ARGV[0].downcase
case command
when "acp"
  message = get_commit_message
  git "add", "."
  git "commit", "-m", message
  git "push"
when "amend"
  message = get_commit_message "new"
  git "commit", "--amend", "--reset", "-m", message
when "last"
  git "log", "-1", "--pretty=format:%h %an: %s (%ad).", "--date=format:%Y-%m-%d %H:%M:%S"
when "log"
  git "log", "--pretty=format:%h %an: %s (%ad).", "--date=format:%Y-%m-%d %H:%M:%S"
when "undo"
  git "reset", "--soft", "HEAD~1"
else
  STDERR.puts "Unknown command: #{command}."
  usage
  exit 1
end
