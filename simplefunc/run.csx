using System.Net;
using System.Diagnostics;
using System;
using System.IO;
using Newtonsoft.Json;
using System.Text;

public static async Task<HttpResponseMessage> Run(HttpRequestMessage req, TraceWriter log)
{
    Process process = new Process();
    process.StartInfo.FileName = "D:/home/site/wwwroot/HttpTriggerCSharp1/main.exe";

    var data = await req.Content.ReadAsStringAsync();
    await WriteToFileAsync(data, log);
    
    process.StartInfo.RedirectStandardOutput = true;
    process.StartInfo.UseShellExecute = false;
    process.Start();
    string q = "";
    while ( ! process.HasExited ) {
        q += process.StandardOutput.ReadToEnd();
    }
    
     log.Info(q);
  
    return req.CreateResponse(HttpStatusCode.OK, "Hello ");
}

static async Task WriteToFileAsync(string text, TraceWriter log)
{
    byte[] buffer = Encoding.UTF8.GetBytes(text);
    Int32 offset = 0;
    Int32 sizeOfBuffer = 4096;
    FileStream fileStream = null;

    fileStream = new FileStream("D:/home/site/wwwroot/HttpTriggerCSharp1/tmp", FileMode.Create, FileAccess.Write, FileShare.None, bufferSize: sizeOfBuffer, useAsync: true);
    await fileStream.WriteAsync(buffer, offset, buffer.Length);
    fileStream.Dispose();
}