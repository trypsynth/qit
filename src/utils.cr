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

  def prompt_single_key(prompt : String) : Char
    print prompt
    STDOUT.flush
    {% if flag?(:win32) %}
      LibC._getch.chr
    {% else %}
      original_termios = uninitialized LibC::Termios
      LibC.tcgetattr(0, pointerof(original_termios))
      raw_termios = original_termios
      raw_termios.c_lflag &= ~(LibC::ICANON | LibC::ECHO)
      LibC.tcsetattr(0, LibC::TCSANOW, pointerof(raw_termios))
      begin
        char = STDIN.read_char
        LibC.tcsetattr(0, LibC::TCSANOW, pointerof(original_termios))
        char || '\0'
      rescue
        LibC.tcsetattr(0, LibC::TCSANOW, pointerof(original_termios))
        '\0'
      end
    {% end %}
  end
end

{% if flag?(:win32) %}
  lib LibC
    fun _getch : Int32
  end
{% else %}
  lib LibC
    struct Termios
      c_iflag : UInt32
      c_oflag : UInt32
      c_cflag : UInt32
      c_lflag : UInt32
      c_cc : UInt8[20]
      c_ispeed : UInt32
      c_ospeed : UInt32
    end

    ICANON  = 2
    ECHO    = 8
    TCSANOW = 0

    fun tcgetattr(fd : Int32, termios : Termios*) : Int32
    fun tcsetattr(fd : Int32, optional_actions : Int32, termios : Termios*) : Int32
  end
{% end %}
