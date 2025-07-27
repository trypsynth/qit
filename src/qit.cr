require "http/client"

USAGE_TEXT = <<-USAGE
Qit - Quin's tiny Git helper.
Usage: qit <command> [<args>...]
Available commands:
  acp <message>: add all files, commit with message, and push.
  amend <message>: amend the last commit with a new message.
  db <name>: delete the local branch <name>.\
  help, -h, --help: show this help message.
  ignore <templates>: download .gitignore template(s) from gitignore.io.
  ignore list: show available templates from gitignore.io.
  last [<number>]: show the last <number> commits (default: 1).
  log: show commit history in readable format.
  nb <name>: switch to branch <name>, creating it if it doesn't exist.
  reset: hard reset to last commit, discarding all changes.
  status: show simplified summary of working directory changes.
  undo: undo last commit while keeping changes intact.
USAGE

def usage
  puts USAGE_TEXT
end

def git(*args : String)
  unless Process.run("git", args, output: STDOUT, error: STDERR).success?
    abort "git #{args.join(" ")} failed."
  end
end

def error_exit(message : String)
  STDERR.puts message
  usage
  exit 1
end

def get_commit_message(prefix : String = "") : String
  return ARGV[1..].join(" ") if ARGV.size >= 2
  desc = prefix.empty? ? "commit message" : "#{prefix} commit message"
  error_exit "Missing #{desc}."
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
when "db"
  branch = ARGV[1]?
  error_exit "Missing branch name." unless branch && !branch.empty?
  delete_branch branch
when "help", "--help", "-h"
  usage
when "ignore"
  case ARGV[1]?.try(&.downcase)
  when "list"
    list_gitignore_templates
  when .nil?
    error_exit "Missing template name(s)."
  else
    download_gitignore ARGV[1]
  end
when "last"
  count = ARGV[1]?.try(&.to_i?) || 1
  git "log", "-#{count}", "--pretty=format:%h %an: %s (%ad).", "--date=format:%Y-%m-%d %H:%M:%S"
when "log"
  git "log", "--pretty=format:%h %an: %s (%ad).", "--date=format:%Y-%m-%d %H:%M:%S"
when "nb"
  branch = ARGV[1]?
  error_exit "Missing branch name." unless branch && !branch.empty?
  switch_or_create_branch branch
when "reset"
  print "This will discard all changes. Continue? (y/N) "
  confirm = gets.try(&.strip.downcase)
  exit unless confirm == "y"
  git "reset", "--hard"
when "status"
  show_status
when "undo"
  git "reset", "--soft", "HEAD~1"
else
  error_exit "Unknown command: #{command}."
end
