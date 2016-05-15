class Exception
  def initialize(*args)
    api = ENV["ERROR_API"] || "localhost:8080"
    params = ENV["ERROR_API_PARAMS"]
    @message = `curl --silent "#{api}/?lang=ruby&full=true#{params}"`
  end

  def message
    @message
  end

  def to_s
    @message
  end

  def backtrace
    # line numbers make it waaay too easy
    []
  end

  def backtrace_locations
    []
  end
end
