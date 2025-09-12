require "../command"

module Qit::Commands
  class NewCommand < Command
    def name : String
      "new"
    end

    def description : String
      "pull and list recent commits"
    end

    def execute(args : Array(String)) : Nil
      old_head = `git rev-parse HEAD`.strip
      Utils.git "pull", quiet: true
      new_head = `git rev-parse HEAD`.strip
      if old_head == new_head
        puts "Nothing new."
      else
        puts "Commits since last pull:"
        Utils.git "log", "#{old_head}..#{new_head}", "--pretty=format:#{Utils::COMMIT_FORMAT}", "--date=format:#{Utils::DATE_FORMAT}"
        puts
      end
    end
  end
end
