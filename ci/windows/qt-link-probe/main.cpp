#include <QApplication>
#include <QJsonObject>
#include <QMediaPlayer>
#include <QNetworkAccessManager>
#include <QPluginLoader>
#include <QStringList>
#include <QVideoWidget>

int main(int argc, char **argv) {
    QApplication app(argc, argv);

    QStringList pluginClasses;
    for (const QStaticPlugin &plugin : QPluginLoader::staticPlugins()) {
        pluginClasses.append(
            plugin.metaData().value(QStringLiteral("className")).toString());
    }
    const QStringList requiredPlugins{
        QStringLiteral("QWindowsIntegrationPlugin"),
        QStringLiteral("QModernWindowsStylePlugin"),
        QStringLiteral("QWindowsMediaPlugin"),
        QStringLiteral("QGifPlugin"),
        QStringLiteral("QICOPlugin"),
        QStringLiteral("QJpegPlugin"),
        QStringLiteral("QSvgPlugin"),
        QStringLiteral("QWebpPlugin"),
        QStringLiteral("QTiffPlugin"),
    };
    for (const QString &plugin : requiredPlugins) {
        if (!pluginClasses.contains(plugin)) {
            return 2;
        }
    }

    QMediaPlayer player;
    QNetworkAccessManager network;
    QVideoWidget video;
    (void)player;
    (void)network;
    (void)video;
    return 0;
}
