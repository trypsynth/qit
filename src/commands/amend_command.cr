require "../command"

module Qit::Commands
  class AmendCommand < Command
    def name : String
      "amend"
    end

    def description : String
      "amend the last commit with a new message"
    end

    def usage : String
      "amend <message>"
    end

    def execute(args : Array(String)) : Nil
      Utils.require_args(args, "Missing new commit message.")
      message = args.join(" ")
      Utils.git "commit", "--amend", "--reset", "-m", message
    end
  end
end
