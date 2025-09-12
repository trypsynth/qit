require "../command"

module Qit::Commands
  class CpCommand < Command
    def name : String
      "cp"
    end

    def description : String
      "commit changes to tracked files with the specified message, and push"
    end

    def execute(args : Array(String)) : Nil
      Utils.require_args(args, "Missing commit message.")
      message = args.join(" ")
      Utils.git "commit", "-am", message
      Utils.git "push"
    end
  end
end
