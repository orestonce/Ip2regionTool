#ifndef DBTOTXTFORM_H
#define DBTOTXTFORM_H

#include <QWidget>

namespace Ui {
class DbToTxtForm;
}

class DbToTxtForm : public QWidget
{
    Q_OBJECT

public:
    explicit DbToTxtForm(QWidget *parent = 0);
    ~DbToTxtForm();
private slots:
    void on_pushButton_selectDb_clicked();

    void on_pushButton_selectTxt_clicked();

    void on_pushButton_startConvert_clicked();
private:
    void refresh_startConvert_Enable();
private:
    Ui::DbToTxtForm *ui;
};

#endif // DBTOTXTFORM_H
