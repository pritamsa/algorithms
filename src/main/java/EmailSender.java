import org.apache.commons.mail.DefaultAuthenticator;
import org.apache.commons.mail.EmailException;
import org.apache.commons.mail.HtmlEmail;
import org.apache.http.client.HttpClient;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.impl.conn.PoolingHttpClientConnectionManager;

public class EmailSender {
    public void method() throws EmailException {
        PoolingHttpClientConnectionManager cm = new PoolingHttpClientConnectionManager();

        cm.setMaxTotal(50);
        cm.setDefaultMaxPerRoute(20);

        CloseableHttpClient cli = HttpClients.custom().setConnectionManager(cm).build();

        HtmlEmail htmlEmail = new HtmlEmail();
        htmlEmail.setHostName("smtp.gmail.com");
        htmlEmail.setSmtpPort(587);
        htmlEmail.setDebug(true);
        htmlEmail.setAuthenticator(new DefaultAuthenticator("userId", "password"));
        htmlEmail.setTLS(true);
        htmlEmail.addTo("recipient@gmail.com", "recipient");
        htmlEmail.setFrom("sender@gmail.com", "sender");
        htmlEmail.setSubject("Send HTML email with body content from URI");
        htmlEmail.setHtmlMsg("msg");
        htmlEmail.send();



    }
}
