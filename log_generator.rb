#!/usr/bin/ruby

class IPGenerator
  public
  def initialize(session_count, session_length)
    @session_count = session_count
    @session_length = session_length

    @sessions = {}
  end

  public
  def get_ip
    session_gc
    session_create

    ip = @sessions.keys[Kernel.rand(@sessions.length)]
    @sessions[ip] += 1
    return ip
  end

  private
  def session_create
    while @sessions.length < @session_count
      @sessions[random_ip] = 0
    end
  end

  private
  def session_gc
    @sessions.each do |ip, count|
      @sessions.delete(ip) if count >= @session_length
    end
  end

  private
  def random_ip
    octets = []
    octets << Kernel.rand(223) + 1
    3.times { octets << Kernel.rand(255) }

    return octets.join(".")
  end
end

class LogGenerator
  EXTENSIONS = {
    'html' => 70,
    'png' => 15,
    'gif' => 10,
    'css' => 5,
  }

  RESPONSE_CODES = {
    200 => 86,
    302 => 6,
    404 => 5,
    503 => 3,
  }

  USER_AGENTS = {
    "Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; Win64; x64; Trident/6.0)" => 12,
    "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)" => 12,
    "Mozilla/5.0 (iPhone; CPU iPhone OS 8_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12B410 Safari/600.1.4" => 12,
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.9; rv:23.0) Gecko/20100101 Firefox/23.0" => 12,
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/29.0.1547.57 Safari/537.36" => 12,
  }

  PATHS = {
    "/articles/bosun/" => 5,
    "/articles/datapower-static-routes/" => 5,
    "/articles/google-code-jam/" => 5,
    "/articles/chess-board-in-objective-c/" => 10,
    "/articles/array-processing-in-ruby/" => 8,
    "/tags/datapower/" => 10,
    "/tags/open-source/" => 10,
    "/tags/ruby/" => 5,
    "/tags/python/" => 10
  }

  FILES = {
    "header" => 5,
    "list" => 4,
    "item" => 3
  }

  public
  def initialize(ipgen)
    @ipgen = ipgen
  end

  public
  def write_qps(dest, qps)
    sleep = 1.0 / qps
    loop do
      write(dest, 1)
      sleep(sleep + Kernel.rand()*0.1)
    end
  end

  public
  def write(dest, count)
    count.times do
      ip = @ipgen.get_ip
      ext = pick_weighted_key(EXTENSIONS)
      resp_code = pick_weighted_key(RESPONSE_CODES)
      resp_size = Kernel.rand(2 * 1024) + 192;
      ua = pick_weighted_key(USER_AGENTS)
      path = pick_weighted_key(PATHS)
      file = pick_weighted_key(FILES)
      date = Time.now.strftime("%d/%b/%Y:%H:%M:%S %z")
      dest.write("#{ip} \"#{ua}\" - [#{date}] \"GET #{path}#{file}.#{ext} HTTP/1.1\" " +
                 "#{resp_code} #{resp_size}\n")
    end
  end

  private
  def pick_weighted_key(hash)
    total = 0
    hash.values.each { |t| total += t }
    random = Kernel.rand(total)

    running = 0
    hash.each do |key, weight|
      if random >= running and random < (running + weight)
        return key
      end
      running += weight
    end

    return hash.keys.first
  end
end

$stdout.sync = true
ipgen = IPGenerator.new(100, 10)
LogGenerator.new(ipgen).write_qps($stdout, 30)
