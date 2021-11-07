#include "mainwindow.h"
#include "ui_mainwindow.h"
#include <QDebug>

MainWindow::MainWindow(QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow)
{
    ui->setupUi(this);
    this->m_form1 = nullptr;
    this->m_form2 = nullptr;
    this->refresh_status(1);
}

MainWindow::~MainWindow()
{
    delete ui;
}

void MainWindow::on_actionDbToTxt_triggered()
{
    this->refresh_status(1);
}

void MainWindow::on_actionTxtToDb_triggered()
{
    this->refresh_status(2);
}

void MainWindow::refresh_status(int idx)
{
    switch (idx) {
    case 1:
        ui->actionDbToTxt->setChecked(true);
        ui->actionTxtToDb->setChecked(false);
        if (this->m_form2 != nullptr) {
            this->m_form2->deleteLater();
            this->m_form2 = nullptr;
        }
        if (this->m_form1 == nullptr) {
            this->m_form1 = new DbToTxtForm(this);
        }
        this->setCentralWidget(this->m_form1);
        break;
    case 2:
        ui->actionDbToTxt->setChecked(false);
        ui->actionTxtToDb->setChecked(true);
        if (this->m_form1 != nullptr) {
            this->m_form1->deleteLater();
            this->m_form1 = nullptr;
        }
        if (this->m_form2 == nullptr) {
            this->m_form2 = new TxtToDbForm(this);
        }
        this->setCentralWidget(this->m_form2);
        break;
    }
}
