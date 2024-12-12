# iOS Library Sample App: TunneledWebRequest

## *\*\* Caveats \*\*\*

### Untunneled OCSP Requests

See "Risk of Online Certificate Status Protocol (OCSP) Leaks" in [../../USAGE.md](../../USAGE.md).

## Configuring, Building, Running

The sample app requires some extra files and configuration before building.

### Get the framework.

1. Run `pod install` 

### Get the configuration.

1. Contact Vaipn Inc. to obtain configuration values to use in your app. 
   (This is requried to use the Vaipn network.)
2. Make a copy of `TunneledWebRequest/vaipn-config.json.stub`, 
   removing the `.stub` extension.
3. Edit `vaipn-config.json`. Remove the comments and fill in the values with 
   those received from Vaipn Inc. The `"ClientVersion"` value is up to you.

### Ready!

TunneledWebRequest should now compile and run.

## License

See the [LICENSE](../LICENSE) file.
