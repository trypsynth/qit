require "../command"

module Qit::Commands
  class AcpCommand < Command
    def name : String
      "acp"
    end

    def description : String
      "add all files, commit with the specified message, and push"
    end

    def usage : String
      "acp <message>"
    end

    def execute(args : Array(String)) : Nil
      Utils.error_exit "Missing commit message." if args.empty?
      message = args.join(" ")
      Utils.git "add", "."
      Utils.git "commit", "-m", message
      Utils.git "push"
    end
  end
end
