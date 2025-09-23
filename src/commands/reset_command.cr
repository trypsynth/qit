require "../command"

module Qit::Commands
  class ResetCommand < Command
    def name : String
      "reset"
    end

    def description : String
      "hard reset to last commit, discarding all changes"
    end

    def usage : String
      "reset"
    end

    def execute(args : Array(String)) : Nil
      confirm = Utils.prompt_single_key("This will discard all changes. Continue? (y/N) ")
      puts
      exit unless confirm.downcase == 'y'
      Utils.git "reset", "--hard"
    end
  end
end
