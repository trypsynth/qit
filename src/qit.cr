require "http/client"

COMMIT_FORMAT = "%h %an: %s (%ad)."
DATE_FORMAT   = "%Y-%m-%d %H:%M:%S"
USAGE_TEXT = <<-USAGE
Qit - Quin's tiny Git helper.
Usage: qit <command> [<args>...]
Available commands:
  acp <message>: add all files, commit with message, and push.
  amend <message>: amend the last commit with a new message.
  cp <message>: commit changes to tracked files with message, and push.
  db <name>: delete the local branch <name>.
  help, -h, --help: show this help message.
  ignore <templates>: download .gitignore template(s) from gitignore.io.
  ignore list: show available templates from gitignore.io.
  last [<number>]: show the last <number> commits (default: 1).
  log: show commit history in readable format.
  nb <name>: switch to branch <name>, creating it if it doesn't exist.
  new: pull and list recent commits.
  reset: hard reset to last commit, discarding all changes.
  status: show simplified summary of working directory changes.
  undo: undo last commit while keeping changes intact.
USAGE

def git(*args : String, quiet : Bool = false)
  output = quiet ? Process::Redirect::Close : STDOUT
  error = quiet ? Process::Redirect::Close : STDERR
  unless Process.run("git", args, output: output, error: error).success?
    abort "git #{args.join(" ")} failed."
  end
end

def error_exit(message : String)
  STDERR.puts message
  exit 1
end

def require_args(args : Array(String), missing_message : String)
  error_exit missing_message if args.empty?
end

def show_help
  puts USAGE_TEXT
end

def get_current_branch : String
  `git rev-parse --abbrev-ref HEAD`.strip
end

def branch_exists?(name : String) : Bool
  !`git branch --list #{name}`.strip.empty?
end

def delete_branch(name : String)
  current = get_current_branch
  error_exit "Cannot delete current branch #{name}." if current == name
  error_exit "Branch '#{name}' does not exist." unless branch_exists?(name)
  git "branch", "-d", name
end

def http_request(url : String, &block : HTTP::Client::Response -> Nil)
  response = HTTP::Client.get(url)
  yield response
rescue ex
  STDERR.puts "Failed to fetch from #{url}: #{ex.message}"
  exit 1
end

def download_gitignore(templates : String)
  url = "https://www.toptal.com/developers/gitignore/api/#{templates}"
  http_request(url) do |response|
    if response.status_code == 200
      File.write(".gitignore", response.body)
      puts "Downloaded .gitignore for: #{templates}"
    else
      STDERR.puts "Error: HTTP #{response.status_code} from gitignore.io"
      exit 1
    end
  end
end

def list_gitignore_templates
  url = "https://www.toptal.com/developers/gitignore/api/list?format=lines"
  http_request(url) do |response|
    if response.status_code == 200
      puts "Available gitignore templates:"
      puts response.body
    else
      STDERR.puts "Error fetching list. HTTP #{response.status_code}"
    end
  end
end

def switch_or_create_branch(name : String)
  current = get_current_branch
  return puts "Already on branch #{name}." if current == name
  if branch_exists?(name)
    puts "Branch #{name} already exists. Switching to it..."
    git "checkout", name
  else
    git "checkout", "-b", name
  end
end

def get_git_status : {Array(String), Array(String)}
  output = IO::Memory.new
  success = Process.run("git", ["status", "--porcelain"], output: output)
  unless success
    STDERR.puts "Failed to get git status."
    exit 1
  end
  staged = [] of String
  unstaged = [] of String
  output.rewind
  output.each_line do |line|
    next if line.strip.empty?
    x, y = line[0]?, line[1]?
    file = line[3..].strip
    staged << file if x != ' ' && x != '?'
    unstaged << file if y != ' '
  end
  {staged, unstaged}
end

def show_status
  staged, unstaged = get_git_status
  return puts "Working tree clean." if staged.empty? && unstaged.empty?
  unless staged.empty?
    puts "Staged for commit:"
    staged.each { |f| puts "  #{f}" }
  end
  unless unstaged.empty?
    puts "Not staged for commit:"
    unstaged.each { |f| puts "  #{f}" }
  end
end

def handle_acp_command(args : Array(String))
  require_args(args, "Missing commit message.")
  message = args.join(" ")
  git "add", "."
  git "commit", "-m", message
  git "push"
end

def handle_amend_command(args : Array(String))
  require_args(args, "Missing new commit message.")
  message = args.join(" ")
  git "commit", "--amend", "--reset", "-m", message
end

def handle_cp_command(args : Array(String))
  require_args(args, "Missing commit message.")
  message = args.join(" ")
  git "commit", "-am", message
  git "push"
end

def handle_db_command(args : Array(String))
  require_args(args, "Missing branch name.")
  delete_branch args[0]
end

def handle_ignore_command(args : Array(String))
  require_args(args, "Missing template name(s). Use 'ignore list' to see available templates.")
  if args[0].downcase == "list"
    list_gitignore_templates
  else
    download_gitignore args[0]
  end
end

def handle_last_command(args : Array(String))
  count = args[0]?.try(&.to_i?) || 1
  git "log", "-#{count}", "--pretty=format:#{COMMIT_FORMAT}", "--date=format:#{DATE_FORMAT}"
end

def handle_nb_command(args : Array(String))
  require_args(args, "Missing branch name.")
  switch_or_create_branch args[0]
end

def handle_new_command(args : Array(String))
  old_head = `git rev-parse HEAD`.strip
  git "pull", quiet: true
  new_head = `git rev-parse HEAD`.strip
  if old_head == new_head
    puts "Nothing new."
  else
    puts "Commits since last pull:"
    git "log", "#{old_head}..#{new_head}", "--pretty=format:#{COMMIT_FORMAT}", "--date=format:#{DATE_FORMAT}"
puts
  end
end

def handle_reset_command(args : Array(String))
  print "This will discard all changes. Continue? (y/N) "
  confirm = gets.try(&.strip.downcase)
  exit unless confirm == "y"
  git "reset", "--hard"
end

if ARGV.empty?
  show_help
  exit 1
end
command = ARGV[0].downcase
args = ARGV[1..]
case command
when "acp"
  handle_acp_command(args)
when "amend"
  handle_amend_command(args)
when "cp"
  handle_cp_command(args)
when "db"
  handle_db_command(args)
when "ignore"
  handle_ignore_command(args)
when "last"
  handle_last_command(args)
when "log"
  git "log", "--pretty=format:#{COMMIT_FORMAT}", "--date=format:#{DATE_FORMAT}"
when "nb"
  handle_nb_command(args)
when "new"
  handle_new_command(args)
when "reset"
  handle_reset_command(args)
when "status"
  show_status
when "undo"
  git "reset", "--soft", "HEAD~1"
when "-h", "--help", "help"
  show_help
else
  error_exit "Unknown command: #{command}."
end
