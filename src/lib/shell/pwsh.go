package shell

import "strings"
func PowershellInit() string {
	return strings.Join([]string{
		UniversalInstaller(),
	}, "\n")
}
func UniversalInstaller() string {
    // powershell function to try install using winget or choco
    return `
		function PowershelUInstaller-Remove {
			param (
				[string[]]$packageNames
			)

			$pkgmgrs = @("winget", "choco")

			foreach ($packageName in $packageNames) {
				if (Get-Command $packageName -ErrorAction SilentlyContinue) {
					Write-Host "$packageName is already installed."
				}
				foreach ($pkgmgr in $pkgmgrs) {
					if (Get-Command $pkgmgr -ErrorAction SilentlyContinue) {
						try {
							Write-Host "Attempting to install $packageName using $pkgmgr..."
							if ($pkgmgr -eq "winget") {
								winget uninstall $packageName -e
							} elseif ($pkgmgr -eq "choco") {
								choco uninstall $packageName -y
							}

							if ($LASTEXITCODE -eq 0) {
								Write-Host "$packageName installed successfully using $pkgmgr."
								break
							}
						} catch {
							Write-Host "Failed to install $packageName using $pkgmgr."
						}
					} else {
						Write-Host "$pkgmgr is not installed."
					}
				}
			}

			Write-Host "Failed to install some packages using all package managers."
		}	
		function PowershelUInstaller {
			param (
				[string[]]$packageNames
			)

			$pkgmgrs = @("winget", "choco")

			foreach ($packageName in $packageNames) {
				if (Get-Command $packageName -ErrorAction SilentlyContinue) {
					Write-Host "$packageName is already installed."
				}
				foreach ($pkgmgr in $pkgmgrs) {
					if (Get-Command $pkgmgr -ErrorAction SilentlyContinue) {
						try {
							Write-Host "Attempting to install $packageName using $pkgmgr..."
							if ($pkgmgr -eq "winget") {
								winget install $packageName -e
							} elseif ($pkgmgr -eq "choco") {
								choco install $packageName -y
							}

							if ($LASTEXITCODE -eq 0) {
								Write-Host "$packageName installed successfully using $pkgmgr."
								break
							}
						} catch {
							Write-Host "Failed to install $packageName using $pkgmgr."
						}
					} else {
						Write-Host "$pkgmgr is not installed."
					}
				}
			}

			Write-Host "Failed to install some packages using all package managers."
		}
	`
}