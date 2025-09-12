require "../command"

module Qit::Commands
  class DbCommand < Command
    def name : String
      "db"
    end

    def description : String
      "delete the specified local branch"
    end

    def execute(args : Array(String)) : Nil
      Utils.require_args(args, "Missing branch name.")
      name = args[0]
      current = Utils.current_branch
      Utils.error_exit "Cannot delete current branch #{name}." if current == name
      Utils.error_exit "Branch '#{name}' does not exist." unless Utils.branch_exists?(name)
      Utils.git "branch", "-d", name
    end
  end
end
