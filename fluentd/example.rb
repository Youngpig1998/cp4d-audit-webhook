def filter(tag, time, record)
  if(File.exist?('/var/log/env'))
    env=Hash[*File.read('/var/log/env').split(/[= \n]+/)]
  else
    env={}
  end

  record = record.merge({:"attachments" => {:"name" => "NA"}})
  record[:"attachments"] = record[:"attachments"].merge({:"content" => {:"sourceCrn" => 10}})
  record[:"attachments"][:"content"] = record[:"attachments"][:"content"].merge({:"kubernetes" => {:"container_id" => env["CONTAINERID"]}})
  record[:"attachments"][:"content"][:"kubernetes"] = record[:"attachments"][:"content"][:"kubernetes"].merge({:"container_name" => env["CONTAINERNAME"]})
  record[:"attachments"][:"content"][:"kubernetes"] = record[:"attachments"][:"content"][:"kubernetes"].merge({:"namespace" => env["NAMESPACE"]})
  record[:"attachments"][:"content"][:"kubernetes"] = record[:"attachments"][:"content"][:"kubernetes"].merge({:"pod" => env["PODIPADDRESS"]})
  record
end

# def code(record)
#   if record.has_key?("key1")
#     record["code"] = record["key1"].to_i
#     record.delete("key1")
#   end
#   record
# end
#
# def message(record)
#   case record["key2"].to_i
#   when 100..200
#     level = "INFO"
#   when 201..300
#     level = "WARN"
#   else
#     level = "ERROR"
#   end
#   record.delete("key2")
#
#   record["message"] = level + ":" + record["key3"]
#   record.delete("key3")
#   record
# end
