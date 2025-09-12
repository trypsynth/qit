require "http/client"

module Utils
  extend self

  COMMIT_FORMAT = "%h %an: %s (%ad)."
  DATE_FORMAT   = "%Y-%m-%d %H:%M:%S"

  def git(*args : String, quiet : Bool = false) : Nil
    output = quiet ? Process::Redirect::Close : STDOUT
    error = quiet ? Process::Redirect::Close : STDERR
    unless Process.run("git", args, output: output, error: error).success?
      abort "git #{args.join(" ")} failed."
    end
  end

  def error_exit(message : String) : NoReturn
    STDERR.puts message
    exit 1
  end

  def require_args(args : Array(String), missing_message : String) : Nil
    error_exit missing_message if args.empty?
  end

  def http_request(url : String, &block : HTTP::Client::Response -> Nil) : Nil
    response = HTTP::Client.get(url)
    yield response
  rescue ex
    STDERR.puts "Failed to fetch from #{url}: #{ex.message}"
    exit 1
  end

  def current_branch : String
    `git rev-parse --abbrev-ref HEAD`.strip
  end

  def branch_exists?(name : String) : Bool
    !`git branch --list #{name}`.strip.empty?
  end
end
