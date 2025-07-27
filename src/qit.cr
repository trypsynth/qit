require "http/client"

def usage
  puts <<-USAGE
  Qit - Quin's tiny Git helper.
  Usage: qit <command> [<args>...]
  Available commands:
    acp <message>: add all files, commit with message, and push.
    amend <message>: amend the last commit with a new message.
    ignore <templates>: download .gitignore template(s) from gitignore.io.
    ignore list: show available templates from gitignore.io.
    last [<number>]: show the last <number> commits (default: 1).
    log: show commit history in readable format.
    reset: hard reset to last commit, discarding all changes.
    status: show simplified summary of working directory changes.
    undo: undo last commit while keeping changes intact.
  USAGE
end

def git(*args : String)
  abort "Command 'git #{args.join(" ")}' failed." unless Process.run("git", args, output: STDOUT, error: STDERR).success?
end

def get_commit_message(prefix : String = "") : String
  if ARGV.size < 2
    desc = prefix.empty? ? "commit message" : "#{prefix} commit message"
    STDERR.puts "Missing #{desc}."
    usage
    exit 1
  end
  ARGV[1..].join(" ")
end

def download_gitignore(templates : String)
  url = "https://www.toptal.com/developers/gitignore/api/#{templates}"
  response = HTTP::Client.get(url)
  if response.status_code == 200
    File.write(".gitignore", response.body)
    puts "Downloaded .gitignore for: #{templates}"
  else
    STDERR.puts "Error: HTTP #{response.status_code} from gitignore.io"
    exit 1
  end
rescue ex
  STDERR.puts "Failed to fetch from gitignore.io: #{ex.message}"
  exit 1
end

def list_gitignore_templates
  url = "https://www.toptal.com/developers/gitignore/api/list?format=lines"
  response = HTTP::Client.get(url)
  if response.status_code == 200
    puts "Available gitignore templates:"
    puts response.body
  else
    STDERR.puts "Error fetching list. HTTP #{response.status_code}"
  end
rescue ex
  STDERR.puts "Error fetching template list: #{ex.message}"
end

def show_status
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
    x = line[0]?
    y = line[1]?
    file = line[3..].strip
    if x != ' ' && x != '?'
      staged << file
    end
    if y != ' '
      unstaged << file
    end
  end
  if staged.empty? && unstaged.empty?
    puts "Working tree clean."
    return
  end
  unless staged.empty?
    puts "Staged for commit:"
    staged.each { |f| puts "  #{f}" }
  end
  unless unstaged.empty?
    puts "Not staged:"
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
when "ignore"
  if ARGV[1]?.try(&.downcase) == "list"
    list_gitignore_templates
    exit
  elsif template = ARGV[1]?
    download_gitignore template
  else
    STDERR.puts "Missing template name(s)."
    usage
    exit 1
  end
when "last"
  count = ARGV[1]?.try(&.to_i?) || 1
  git "log", "-#{count}", "--pretty=format:%h %an: %s (%ad).", "--date=format:%Y-%m-%d %H:%M:%S"
when "log"
  git "log", "--pretty=format:%h %an: %s (%ad).", "--date=format:%Y-%m-%d %H:%M:%S"
when "reset"
  git "reset", "--hard"
when "status"
  show_status
when "undo"
  git "reset", "--soft", "HEAD~1"
else
  STDERR.puts "Unknown command: #{command}."
  usage
  exit 1
end
