require "../command"
require "../utils"

module Qit::Commands
  class UndoCommand < Command
    def name : String
      "undo"
    end

    def description : String
      "undo last commit while keeping changes intact"
    end

    def execute(args : Array(String)) : Nil
      Utils.git "reset", "--soft", "HEAD~1"
    end
  end
end
