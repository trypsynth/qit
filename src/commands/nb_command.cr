require "../command"

module Qit::Commands
  class NbCommand < Command
    def name : String
      "nb"
    end

    def description : String
      "switch to the branch with  the specified name, creating it if it doesn't exist"
    end

    def usage : String
      "nb <branch_name>"
    end

    def execute(args : Array(String)) : Nil
      Utils.require_args(args, "Missing branch name.")
      name = args[0]
      current = Utils.current_branch
      return puts "Already on branch #{name}." if current == name
      if Utils.branch_exists?(name)
        puts "Branch #{name} already exists. Switching to it..."
        Utils.git "checkout", name
      else
        Utils.git "checkout", "-b", name
      end
    end
  end
end
