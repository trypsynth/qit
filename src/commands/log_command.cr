require "../command"

module Qit::Commands
  class LogCommand < Command
    def name : String
      "log"
    end

    def description : String
      "show commit history in readable format"
    end

    def usage : String
      "log"
    end

    def execute(args : Array(String)) : Nil
      Utils.git "log", "--pretty=format:#{Utils::COMMIT_FORMAT}", "--date=format:#{Utils::DATE_FORMAT}"
      puts
    end
  end
end
