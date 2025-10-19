# File: Measure-Memory.ps1
function Measure-Memory {
    param(
        [Parameter(Mandatory = $true, ValueFromRemainingArguments = $true)]
        [string[]]$Command
    )

    $start = Get-Date

    $exe = $Command[0]
    $args = if ($Command.Length -gt 1) { $Command[1..($Command.Length-1)] } else { @() }

    $psi = New-Object System.Diagnostics.ProcessStartInfo
    $psi.FileName = $exe
    $psi.Arguments = ($args -join ' ')
    $psi.RedirectStandardOutput = $true
    $psi.RedirectStandardError  = $true
    $psi.UseShellExecute = $false
    $psi.CreateNoWindow = $true

    $proc = New-Object System.Diagnostics.Process
    $proc.StartInfo = $psi
    $proc.Start() | Out-Null

    $maxMemory = 0

    # Real-time output/error
    while (-not $proc.HasExited) {
        while (-not $proc.StandardOutput.EndOfStream) {
            $line = $proc.StandardOutput.ReadLine()
            if ($line) { Write-Host $line }
        }
        while (-not $proc.StandardError.EndOfStream) {
            $line = $proc.StandardError.ReadLine()
            if ($line) { Write-Host $line -ForegroundColor Red }
        }

        try {
            $mem = (Get-Process -Id $proc.Id -ErrorAction Stop).WorkingSet64 / 1MB
            if ($mem -gt $maxMemory) { $maxMemory = $mem }
        } catch {}
        Start-Sleep -Milliseconds 100
    }

    # Remaining buffered lines
    while (-not $proc.StandardOutput.EndOfStream) {
        $line = $proc.StandardOutput.ReadLine()
        if ($line) { Write-Host $line }
    }
    while (-not $proc.StandardError.EndOfStream) {
        $line = $proc.StandardError.ReadLine()
        if ($line) { Write-Host $line -ForegroundColor Red }
    }

    $proc.WaitForExit()
    $elapsed = (Get-Date) - $start

    Write-Host ""
    Write-Host ("Exit code: {0}" -f $proc.ExitCode)
    Write-Host ("Elapsed:   {0:N2} s" -f $elapsed.TotalSeconds)
    Write-Host ("Peak RAM:  {0:N2} MB" -f $maxMemory)
}
