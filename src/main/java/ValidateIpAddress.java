import org.apache.commons.validator.routines.InetAddressValidator;

public class ValidateIpAddress {

    public boolean isValid(String ip) {
        InetAddressValidator validator = InetAddressValidator.getInstance();

        if (validator.isValid(ip)) {

        }

        return false;

    }
}
