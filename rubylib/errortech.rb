class Exception
  def initialize(*args)
    api = ENV["ERROR_API"] || "localhost:8080"
    @message = `curl --silent "#{api}/?lang=ruby&full=true"`
  end

  def message
    @message
  end

  def to_s
    @message
  end

  def backtrace
    []
  end
end
