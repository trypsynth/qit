require "../command"

module Qit::Commands
  class ResetCommand < Command
    def name : String
      "reset"
    end

    def description : String
      "hard reset to last commit, discarding all changes"
    end

    def execute(args : Array(String)) : Nil
      print "This will discard all changes. Continue? (y/N) "
      confirm = gets.try(&.strip.downcase)
      exit unless confirm == "y"
      Utils.git "reset", "--hard"
    end
  end
end
