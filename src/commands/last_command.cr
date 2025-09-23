require "../command"

module Qit::Commands
  class LastCommand < Command
    def name : String
      "last"
    end

    def description : String
      "show the last <number> commits (default: 1)"
    end

    def usage : String
      "last [<number>]"
    end

    def execute(args : Array(String)) : Nil
      count = args[0]?.try(&.to_i?) || 1
      Utils.git "log", "-#{count}", "--pretty=format:#{Utils::COMMIT_FORMAT}", "--date=format:#{Utils::DATE_FORMAT}"
      puts
    end
  end
end
