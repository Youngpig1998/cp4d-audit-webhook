def filter(tag, time, record)
  if(File.exist?('/var/log/env'))
    env=Hash[*File.read('/var/log/env').split(/[= \n]+/)]
  else
    env={}
  end
  $logSourceCRN = record[:"logSourceCRN"]
  $message = record[:"message"]

  record = record.merge({:"attachments" => [{:"name" => "ibm-cp-data"}]})
  record[:"attachments"][0] = record[:"attachments"][0].merge({:"contentType" => "http://schemas.ibm.com/cloud/content/1.0/cloudpak"})
  record[:"attachments"][0] = record[:"attachments"][0].merge({:"logSourceCRN" => $logSourceCRN})
  record[:"attachments"][0] = record[:"attachments"][0].merge({:"content" => {:"sourceCrn" => 10}})

  record[:"attachments"][0][:"content"] = record[:"attachments"][0][:"content"].merge({:"message" => $message})
  record[:"attachments"][0][:"content"] = record[:"attachments"][0][:"content"].merge({:"kubernetes" => {:"container_id" => env["CONTAINERID"]}})
  record[:"attachments"][0][:"content"][:"kubernetes"] = record[:"attachments"][0][:"content"][:"kubernetes"].merge({:"container_name" => env["CONTAINERNAME"]})
  record[:"attachments"][0][:"content"][:"kubernetes"] = record[:"attachments"][0][:"content"][:"kubernetes"].merge({:"namespace" => env["NAMESPACE"]})
  record[:"attachments"][0][:"content"][:"kubernetes"] = record[:"attachments"][0][:"content"][:"kubernetes"].merge({:"pod" => env["PODIPADDRESS"]})
  record
end