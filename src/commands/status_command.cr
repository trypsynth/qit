require "../command"

module Qit::Commands
  class StatusCommand < Command
    def name : String
      "status"
    end

    def description : String
      "show simplified summary of working directory changes"
    end

    def usage : String
      "status"
    end

    def execute(args : Array(String)) : Nil
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

    private def get_git_status : {Array(String), Array(String)}
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
  end
end
